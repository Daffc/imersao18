import { Controller, Get, Post, Body, Patch, Param, Delete, Query, HttpCode } from '@nestjs/common';
import { SpotsService } from '@app/core/spots/spots.service';
import { CriarLugarRequest } from './request/criar-lugar.request';
import { AtualizarLugarRequest } from './request/atualizar-lugar.request';

@Controller('eventos/:eventId/lugares')
export class LugaresController {
  constructor(private readonly spotsService: SpotsService) { }

  @Post()
  create(@Param('eventId') eventId: string, @Body() criarLugarRequest: CriarLugarRequest) {
    return this.spotsService.create({
      eventId,
      name: criarLugarRequest.nome
    });
  }

  @Get()
  findAll(@Param('eventId') eventId: string) {
    return this.spotsService.findAll(eventId);
  }

  @Get(':spotId')
  findOne(@Param('spotId') spotId: string, @Param('eventId') eventId: string) {
    return this.spotsService.findOne(eventId, spotId);
  }

  @Patch(':spotId')
  update(@Param('spotId') spotId: string, @Param('eventId') eventId: string, @Body() atualizarLugarRequest: AtualizarLugarRequest) {
    return this.spotsService.update(
      eventId,
      spotId,
      {
        name: atualizarLugarRequest.nome
      }
    );
  }

  @HttpCode(204)
  @Delete(':spotId')
  remove(@Param('spotId') spotId: string, @Param('eventId') eventId: string) {
    return this.spotsService.remove(eventId, spotId);
  }
}
